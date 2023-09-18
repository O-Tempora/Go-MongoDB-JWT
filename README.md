How to run:
```
make run
```

How to stop: 
```
ctrl+c THEN make killdb
```

How to restart with new db:
```
make reset
```

Testing process:
```
make run
docker exec -it mongotest mongosh
use gojwt
db.users.find({})
Copy user _id (without brackets and quotes)
Put it as a query param named "guid" in access route
```

To use swagger:
```
http://localhost:5005/swagger/
```