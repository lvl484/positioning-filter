## Positioning filter service (FS)

This service is intended to match positioning messages to a set of filtering rules and pereat messages to a topic in case of they are matched.


#### Position message model

```
    Position {
        user         // id of user that produced position notification event
        latitude     // obviously latitude component of coordinates
        longitude    // obviously longitude component of coordinates
        timestamp    // time when position collected
        arrival      // time when position accepted by system
    }

```

Position messages will be read from a kafka topic named 'positions'. Each position message will come through a set of filters that are related to the user that generated the message. If position is inside the area of the filter it will be forwarded to another kafka topic called 'matched-positions', then notification service will consume them from the topic and push notification by destination.

#### Filter model

Each user will have a set of filters that will be applied to each message produced by the user. If filter is reverted - rule will be matched if position is outside the area. Minimun required zize of the area is 400m2 maximum is 5km2

```
    Filter {
        type           // type of area matching condition
        configuration  // object with configuration fitting the type
        reversed       // flag that describes if rule should match when does not satisfy the condition
        user id        // id of the user that configured the filter
    }
```

###### Filter types

Filter type can be rectangular or round(in future maybe shaped)

**Round filter:** detects if the position is inside of round area

```
    RoundFilter {
        center coordinates {latitude, longitude}    // A point on a map that is the center of the square area
        radius {meters}                             // The radius of the area
    }
```

**Rectangular filter:** detects the position inside of rectangular area

```
    RectangularFilter {
        point top left, point bottom right {latitude, longitude}   // A couple of points on a map that describes margin of rectangle area.
    }
```  

## MIGRATIONS

Migrate with using golang/migrate default table/rules for user
In table must be created default user with name 'fuser' without any rules with own password

### environment variables

DBPORT -host port for database  
POSTGRES_USER -default owner for database  
POSTGRES_PASSWORD -password for owner  
POSTGRES_DB -default database name (db postgres will be also created)  
MIGRATIONS_PATH -path to folder with .sql migrations files  

### commands

docker run -v "$MIGRATIONS_PATH:/migrations"  --network host migrate/migrate -path=/migrations/ -database "postgresql://$POSTGRES_USER:$POSTGRES_PASSWORD@localhost:$DB_POST/$POSTGRES_DB?sslmode=disable" up
