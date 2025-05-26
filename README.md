# Recipe Video Parser
A Golang application that parses recipe videos from various platforms and extracts the recipe details. Tested with YouTube and TikTok.

## Installation
1. Install gallery-dl and yt-dlp:
```bash
python3 -m pip install -U gallery-dl yt-dlp
```
2. Clone the repository

3. Install Go dependencies:
```bash
make install
```
4. Set the environment variables using `example.env` as a template

5. Build the application:
```bash
make build
```
6. Run the application:
```bash
./bin/api
```

## TODO
- [ ] Add support for platforms that require authentication.
- [ ] Implement a web interface for easier access.
