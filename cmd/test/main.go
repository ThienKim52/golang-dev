package main

import (

	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog"
	
)

func main() {
	// rclient, err := redis.NewClient("")
	// if err != nil {
	// 	panic(err)
	// }
	// rclient.Set(context.Background(), "1235", "google.com", time.Hour)

	// rclient2, err := redis.NewClient("CACHE")
	// if err != nil {
	// 	panic(err)
	// }
	// rclient2.Set(context.Background(), "1999", "google.com", time.Hour)
	zerolog.SetGlobalLevel(zerolog.WarnLevel)
	log.Debug().Msg("test")
	log.Info().Msg("test")
	log.Warn().Msg("test")
	log.Error().Msg("test")
}
