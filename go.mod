module github.com/gravity-corp/balancer

go 1.18

retract (
    v1.0.6 // I tried to delete this repo from pkg.go.dev,
    v1.0.5 // but retract work over way and as you can see
    v1.0.4 // it still exists and available,
    v1.0.3 // despite switching visibility on github.
    v1.0.2 // So I made a mistake that I have my worst
    v1.0.1 // nightmares about and I will never be able
    v1.0.0 // to forget about it. Take care of youself 
    v0.9.9 // and be carefull next time you "go get".
)
