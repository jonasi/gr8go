module github.com/jonasi/gr8go/http

go 1.15

require (
	github.com/jonasi/gr8go/log v0.0.0-00010101000000-000000000000
	github.com/jonasi/gr8go/service v0.0.0-20201224203359-ea1347dfdd1c
	github.com/julienschmidt/httprouter v1.3.0
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822
)

replace (
	github.com/jonasi/gr8go/log => ../log
	github.com/jonasi/gr8go/service => ../service
)
