
#!/bin/bash
curl -d "reset" http://localhost:17000

coordX1=10
coordY1=10
coordX2=790 
coordY2=790 

curl -d "green" http://localhost:17000 
curl -d "bgrect $coordX1 $coordY1 $coordX2 $coordY2 $1" http://localhost:17000 
curl -d "update" http://localhost:17000 