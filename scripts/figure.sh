#!/bin/bash
curl -d "reset" http://localhost:17000

curl -d "white" http://localhost:17000 
curl -d "bgrect 200 200 600 600" http://localhost:17000 
curl -d "figure 400 400" http://localhost:17000 
curl -d "green" http://localhost:17000 
curl -d "figure 480 480" http://localhost:17000 
curl -d "update" http://localhost:17000 