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

## Start all services  

You must create .env file with following variables in the root of project  

### Environment variables  

POSTGRES_USER -default owner for database  
POSTGRES_PASSWORD -password for owner  
POSTGRES_DB -default database name (db postgres will be also created)  
PF_USER -default user for accesing database from api  
PF_PASSWORD -password for user

### Commands  

```bash
docker-compose up --build -d  
```

### Known issues  
  
Matching position with round filters around twelve meridian can returns wrong result  

# examples  
  
```
with (positionLatitude = 179, filterCenterLatitude = 178, filterRadius = 10) matchRound returns true  
with (positionLatitude = -179, filterCenterLatitude = 178, filterRadius = 10) matchRound returns false  
if filter contains area from two sides of twelve meridian and
position and filter center are on different sides, matchRound returns wrong result
```  
