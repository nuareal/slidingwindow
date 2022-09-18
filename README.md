# slidingwindow
An implementation of the sliding window algorithm for counting number of received http requests.
The algorithm is implemented with circular arrays. 

# Features
- Endpoint responds with number of requests received in last minute.
- Data is persistent. When the application is terminated, the data is saved to a cache file and loaded on restart.
- Only uses standard go packages.

# Usage
- Run main.go
- Curl localhost:8080

# Future Improvements
- Add versatility (choose the cache file name, port number, time domain(minute, hour, day))
- Add unit tests
