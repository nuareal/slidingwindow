# slidingwindow
An implementation of the sliding window algorithm for counting number of received http requests.
The algorithm is implemented with circular arrays. 

# Features
- Endpoint responds with number of requests received in last minute.
- Data is persistent. When the application is terminated, the data is saved to a cache file and loaded on restart.
- Only uses standard go packages.

# Future Improvements
- Add versatility (choose the cache file name, time domain(minute, hour, day))
- Add unit tests
