all: run

run:
	go run main.go itog.go lig.go result.go calcule.go

push:
	git push git@github.com:RB-PRO/soccer365.git

pull:
	git pull git@github.com:RB-PRO/soccer365.git

