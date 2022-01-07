# Go Socket Chat

Small socket-based chat app written in Go, HTML, JS, and TailwindCSS

To run locally:

Frontend precompilation:
```
cd static/
npm i
npx tailwind -i ./css/main.css -o ./css/out.css --watch
```

Run app from main directory:
```
go run cmd/web/*.go
```
