# Recipe Video Parser
A Golang application that parses recipe videos from various platforms and extracts the recipe details. Tested with YouTube, TikTok, and Instagram.

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

## Usage

### CLI
- Extract recipe from a URL:
  ```bash
  ./bin/cli extract-recipe <video_url>
  ```
  Extracts the recipe in JSON format from the given URL (YouTube, TikTok, Instagram, etc).

- Create user:
  ```bash
  ./bin/cli create-user <username>
  ```
  Creates a user and generates an API key (only needed for HTTP API access).

- Update user's API key:
  ```bash
  ./bin/cli update-api-key <username>
  ```
  Generates a new API key for the user.

- Get user by username:
  ```bash
  ./bin/cli get-user <username>
  ```
  Shows user data by ID.

- Get extraction summary by user:
  ```bash
  ./bin/cli get-user-summary <username>
  ```
  Shows a monthly summary of extractions and tokens used.

### API

To start the API:
```bash
./bin/api
```

The API requires authentication via API key. You must create users and obtain their API keys using the CLI before you can access the protected endpoint (`/recipes/extract`).

**Authentication:**
All API requests must include the API key in the `Authorization` header using the Bearer scheme:

```
Authorization: Bearer <API_KEY>
```

## Videos with login requirements
For platforms that require login (like Instagram), you can specify a custom `gallery-dl` configuration file in the `.env` file:

```env
GALLERY_CONFIGFILE=/path/to/your/gallery-dl.conf
```

Sample `gallery-dl.conf` for Instagram:
```json
{
		"extractor": {
				"instagram": {
						"cookies": "/path/to/your/instagram_cookies.txt",
				}
		}
}
```

You can obtain the cookies file using browser extensions. More information on how to obtain cookies can be found in the [gallery-dl documentation](https://github.com/mikf/gallery-dl?tab=readme-ov-file#cookies).

More information on how to configure `gallery-dl` can be found in the [gallery-dl documentation](https://gdl-org.github.io/docs/configuration.html)

## TODO
- [ ] Implement a web interface for easier access.
