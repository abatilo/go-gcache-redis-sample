package cachedsvc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/abatilo/go-kube-shutdown/pkg/shutdown"
	"github.com/bluele/gcache"
	"github.com/go-redis/redis/v7"
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/cobra"
)

var (
	// Cmd is the exported cobra command which starts the webhook handler service
	Cmd = &cobra.Command{
		Use:   "svc",
		Short: "Runs the web service",
		Run: func(cmd *cobra.Command, args []string) {
			main()
		},
	}
)

type server struct {
	router      *httprouter.Router
	cache       gcache.Cache
	redisClient *redis.Client
}

func (s *server) setupHandlers() {
	if s.router == nil {
		s.router = httprouter.New()
	}
	s.router.HandlerFunc("GET", "/:key", func(w http.ResponseWriter, r *http.Request) {
		params := httprouter.ParamsFromContext(r.Context())
		key := params.ByName("key")

		val, err := s.cache.Get(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		m := make(map[string]string)
		keysRaw := s.cache.Keys(false)
		for _, kRaw := range keysRaw {
			k := kRaw.(string)
			val, _ := s.cache.Get(k)
			m[k] = val.(string)
		}

		fmt.Println(s.cache.HitRate())
		w.Write([]byte(val.(string)))
	})

	s.router.HandlerFunc("POST", "/:key", func(w http.ResponseWriter, r *http.Request) {
		kv := make(map[string]string)
		err := json.NewDecoder(r.Body).Decode(&kv)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Printf("%v\n", kv)
	})
}

func (s *server) newHTTPServer() *http.Server {
	s.setupHandlers()
	return &http.Server{
		Addr:    ":8000",
		Handler: s.router,
	}
}

func main() {
	client := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName: "master",
		SentinelAddrs: []string{
			"redis-sentinel1:26379",
		},
	})

	gc := gcache.New(50).
		LFU().
		LoaderFunc(func(key interface{}) (interface{}, error) {
			s := key.(string)

			rev, err := client.Get(s).Result()
			if err != redis.Nil {
				// If Redis has it, just send it back
				return rev, nil
			}

			// Reverse the string
			runes := []rune(s)
			for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
				runes[i], runes[j] = runes[j], runes[i]
			}
			res := string(runes)

			// Set the key in Redis
			err = client.Set(s, res, time.Minute*5).Err()
			if err != nil {
				fmt.Printf("Err: %v\n", err)
			}
			return res, err
		}).
		Build()

	srv := &server{cache: gc, redisClient: client}

	fmt.Println("Starting server...")
	err := shutdown.StartSafeServer(srv.newHTTPServer(), "/tmp/live")
	if err != http.ErrServerClosed {
		fmt.Printf("Server didn't shutdown cleanly: %v\n", err)
	}
}
