package main

import (
	"log"
	"net/http"
	"text/template"

	"github.com/anthonyP1202/Hangman"
)

// type HangManData struct {
// 	Word             []rune
// 	WordToFind       []string
// 	GivenLetter      []string
// 	GivenWord        []string
// 	NbrOfAttempt     int
// 	VictoryCondition int
// }

type user struct {
	Name          string
	HighestStreak int
	CurrentStreak int
	TotalWin      int
	TotalPlay     int
}

func Home(w http.ResponseWriter, r *http.Request, user *user) {
	template, err := template.ParseFiles("./index.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, nil)
}

func Info(w http.ResponseWriter, r *http.Request, info *Hangman.HangManData, user *user) {
	template, err := template.ParseFiles("./page/info.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, info)
}

func Final(w http.ResponseWriter, r *http.Request, info *Hangman.HangManData, user *user) {

	template, err := template.ParseFiles("./page/final.html")
	print(user.Name)
	if r.FormValue("name") != "" {
		if user.Name == "" {
			nameChange(user, r.FormValue("name"))
		} else {
			user.Name = r.FormValue("name")
			user.TotalPlay = 0
			user.TotalWin = 0
		}
	}

	if err != nil {
		log.Fatal(err)
	}
	
	if r.FormValue("select") != "" {
		Hangman.HangmanADV(info, r.FormValue("select"))
		user.TotalPlay++
	}

	given := []rune(r.FormValue("letter"))

	for i := 0; i < len(given); i++ {
		if given[i] > 'A' && given[i] < 'Z' {
			given[i] = given[i] + 32
		}
	}

	if given != nil {
		if len(given) >= 1 {
			Hangman.Playingadv(info, given)
		}
	}

	if info.NbrOfAttempt <= 0 || info.VictoryCondition == 1 {
		for i := 0; i < len(info.Word); i++ {
			info.WordToFind[i] = string(info.Word[i])
		}
		if info.VictoryCondition == 1 {
			user.TotalWin++
			user.CurrentStreak++
			if user.CurrentStreak > user.HighestStreak {
				user.HighestStreak = user.CurrentStreak
			}
		}
		http.Redirect(w, r, "win", http.StatusSeeOther)
	}

	template.Execute(w, info)

}

func Win(w http.ResponseWriter, r *http.Request, info *Hangman.HangManData) {
	template, err := template.ParseFiles("./page/win.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, info)
}

func main() {
	var info Hangman.HangManData

	var user user

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Home(w, r, &user)
	})

	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		Info(w, r, &info, &user)
	})

	http.HandleFunc("/final", func(w http.ResponseWriter, r *http.Request) {
		Final(w, r, &info, &user)
	})

	http.HandleFunc("/win", func(w http.ResponseWriter, r *http.Request) {
		Win(w, r, &info)
	})
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static", fs))
	http.ListenAndServe(":8080", nil)
}

func nameChange(user *user, newName string) *user {
	user.Name = newName
	return user
}
