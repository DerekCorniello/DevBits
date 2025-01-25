# Connection and Deployment Instructions

1. `cd PROJECT_ROOT/DevBits/deployment`
2. `chmod +x connect.sh`
3. `./connect.sh`
4. `cd DevBits/backend`
5. `nohup go run your-server.go > server.log 2>&1 &`
6. Close the terminal and the app is live!
