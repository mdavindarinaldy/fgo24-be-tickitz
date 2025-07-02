package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Movie struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	ReleaseDate string `json:"release_date"`
	Poster      string `json:"poster_path"`
	Backdrop    string `json:"backdrop_path"`
}

type MovieDetail struct {
	Overview string `json:"overview"`
	Runtime  int    `json:"runtime"`
}

type MovieDirectors struct {
	Credits Crews `json:"credits"`
}

type Crews struct {
	Crew []Crew `json:"crew"`
}

type Crew struct {
	Name string `json:"name"`
	Job  string `json:"job"`
}

type MovieCasts struct {
	Credits Casts `json:"credits"`
}

type Casts struct {
	Cast []Cast `json:"cast"`
}

type Cast struct {
	Name string `json:"name"`
}

type NowPlayingResponse struct {
	Results []Movie `json:"results"`
	Page    int     `json:"page"`
}

func fetchNowPlayingMovies(apiKey string) error {
	url := "https://api.themoviedb.org/3/movie/now_playing?language=en-US&page=1"

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("gagal membuat request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("gagal melakukan request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request gagal dengan status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("gagal membaca respons: %v", err)
	}

	var nowPlaying NowPlayingResponse
	if err := json.Unmarshal(body, &nowPlaying); err != nil {
		return fmt.Errorf("gagal parsing JSON: %v", err)
	}

	// fmt.Println(nowPlaying)
	// InsertMovies(nowPlaying)
	// InsertDirectors(nowPlaying)
	// InsertCast(nowPlaying)
	return nil
}

func FetchDetailMovie(id int) (MovieDetail, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%d?&append_to_response=credits", id)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return MovieDetail{}, fmt.Errorf("gagal membuat request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiIzZWMxNTY1OTEwMGQ3M2RlNTI0YTlkN2I4NGUxYWE0OCIsIm5iZiI6MTc0NzUzMDg1MC4xMDA5OTk4LCJzdWIiOiI2ODI5MzQ2MmE5MjZlNzQ1Njk1YjUyZTEiLCJzY29wZXMiOlsiYXBpX3JlYWQiXSwidmVyc2lvbiI6MX0.1WmayjjDNmTjMaAUPyPvM9LDJ2GTzrF3rVyYdhLN28k"))
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return MovieDetail{}, fmt.Errorf("gagal melakukan request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return MovieDetail{}, fmt.Errorf("request gagal dengan status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return MovieDetail{}, fmt.Errorf("gagal membaca respons: %v", err)
	}

	var movieDetail MovieDetail
	if err := json.Unmarshal(body, &movieDetail); err != nil {
		return MovieDetail{}, fmt.Errorf("gagal parsing JSON: %v", err)
	}

	return movieDetail, nil
}

func FetchDirectors(id int) (MovieDirectors, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%d?&append_to_response=credits", id)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return MovieDirectors{}, fmt.Errorf("gagal membuat request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiIzZWMxNTY1OTEwMGQ3M2RlNTI0YTlkN2I4NGUxYWE0OCIsIm5iZiI6MTc0NzUzMDg1MC4xMDA5OTk4LCJzdWIiOiI2ODI5MzQ2MmE5MjZlNzQ1Njk1YjUyZTEiLCJzY29wZXMiOlsiYXBpX3JlYWQiXSwidmVyc2lvbiI6MX0.1WmayjjDNmTjMaAUPyPvM9LDJ2GTzrF3rVyYdhLN28k"))
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return MovieDirectors{}, fmt.Errorf("gagal melakukan request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return MovieDirectors{}, fmt.Errorf("request gagal dengan status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return MovieDirectors{}, fmt.Errorf("gagal membaca respons: %v", err)
	}

	var movieDetail MovieDirectors
	if err := json.Unmarshal(body, &movieDetail); err != nil {
		return MovieDirectors{}, fmt.Errorf("gagal parsing JSON: %v", err)
	}

	return movieDetail, nil
}

func FetchCasts(id int) (MovieCasts, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%d?&append_to_response=credits", id)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return MovieCasts{}, fmt.Errorf("gagal membuat request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiIzZWMxNTY1OTEwMGQ3M2RlNTI0YTlkN2I4NGUxYWE0OCIsIm5iZiI6MTc0NzUzMDg1MC4xMDA5OTk4LCJzdWIiOiI2ODI5MzQ2MmE5MjZlNzQ1Njk1YjUyZTEiLCJzY29wZXMiOlsiYXBpX3JlYWQiXSwidmVyc2lvbiI6MX0.1WmayjjDNmTjMaAUPyPvM9LDJ2GTzrF3rVyYdhLN28k"))
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return MovieCasts{}, fmt.Errorf("gagal melakukan request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return MovieCasts{}, fmt.Errorf("request gagal dengan status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return MovieCasts{}, fmt.Errorf("gagal membaca respons: %v", err)
	}

	var movieDetail MovieCasts
	if err := json.Unmarshal(body, &movieDetail); err != nil {
		return MovieCasts{}, fmt.Errorf("gagal parsing JSON: %v", err)
	}

	return movieDetail, nil
}

func Fetch() {
	apiKey := "eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiIzZWMxNTY1OTEwMGQ3M2RlNTI0YTlkN2I4NGUxYWE0OCIsIm5iZiI6MTc0NzUzMDg1MC4xMDA5OTk4LCJzdWIiOiI2ODI5MzQ2MmE5MjZlNzQ1Njk1YjUyZTEiLCJzY29wZXMiOlsiYXBpX3JlYWQiXSwidmVyc2lvbiI6MX0.1WmayjjDNmTjMaAUPyPvM9LDJ2GTzrF3rVyYdhLN28k"

	if err := fetchNowPlayingMovies(apiKey); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	// InsertGenre()
	// InsertMoviesCasts()
}

func InsertMovies(npr NowPlayingResponse) error {
	conn, err := DBConnect()
	if err != nil {
		return err
	}
	defer conn.Close()

	for _, v := range npr.Results {
		detail, _ := FetchDetailMovie(v.Id)
		conn.Exec(
			context.Background(),
			`
			INSERT INTO movies (created_by, title, synopsis, release_date, price, runtime, poster, backdrop, created_at)
			VALUES
			($1,$2,$3,$4,$5,$6,$7,$8,$9);
			`,
			1, v.Title, detail.Overview, v.ReleaseDate, 15000, detail.Runtime, v.Poster, v.Backdrop, time.Now())
	}
	return nil
}

func InsertDirectors(npr NowPlayingResponse) error {
	conn, err := DBConnect()
	if err != nil {
		return err
	}
	defer conn.Close()
	for _, v := range npr.Results {
		detail, _ := FetchDirectors(v.Id)
		fmt.Print(detail)
		directors := []Crew{}
		for _, i := range detail.Credits.Crew {
			if i.Job == "Director" {
				directors = append(directors, i)
			}
		}
		for _, j := range directors {
			conn.Exec(
				context.Background(),
				`
				INSERT INTO directors (name, created_at)
				VALUES
				($1,$2);
				`,
				j.Name, time.Now())
		}
	}
	return nil
}

func InsertCast(npr NowPlayingResponse) error {
	conn, err := DBConnect()
	if err != nil {
		return err
	}
	defer conn.Close()
	for _, v := range npr.Results {
		detail, _ := FetchCasts(v.Id)
		fmt.Print(detail)
		for i := 0; i < 5; i++ {
			conn.Exec(
				context.Background(),
				`
				INSERT INTO casts (name, created_at)
				VALUES
				($1,$2);
				`,
				detail.Credits.Cast[i].Name, time.Now())
		}
	}
	return nil
}

func InsertGenre() error {
	Genres := []string{
		"Action",
		"Adventure",
		"Animation",
		"Comedy",
		"Crime",
		"Documentary",
		"Drama",
		"Family",
		"Fantasy",
		"History",
		"Horror",
		"Music",
		"Mystery",
		"Romance",
		"Sci-Fi",
		"TV Movie",
		"Thriller",
		"War",
		"Western",
	}
	conn, err := DBConnect()
	if err != nil {
		return err
	}
	defer conn.Close()
	for _, j := range Genres {
		conn.Exec(
			context.Background(),
			`
			INSERT INTO genres (name, created_at)
			VALUES
			($1,$2);
			`,
			j, time.Now())
	}
	return nil
}

func InsertMoviesCasts() error {
	conn, err := DBConnect()
	if err != nil {
		return err
	}
	defer conn.Close()
	cast := 1
	for i := 1; i <= 20; i++ {
		for j := 1; j <= 5; j++ {
			conn.Exec(
				context.Background(),
				`
				INSERT INTO movies_casts (id_cast, id_movie)
				VALUES
				($1,$2);
				`, cast, i)
			cast++
		}
	}
	return nil
}
