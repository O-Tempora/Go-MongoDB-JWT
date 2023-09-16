How to run:
    make run

How to stop: 
    ctrl+c THEN make killdb

How to restart with new db:
    make reset

Testing process:
    1. make run
    2. docker exec -it mongotest mongosh
    3. use gojwt
    4. db.users.find({})
    5. Copy user _id (without brackets and quotes)
    6. Put it as a query param named "guid" in access route