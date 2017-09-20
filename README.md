# uploader

Repo uploader holds a minimal service for accepting uploads over HTTP.

## Testing

Using `curl`, it's easy to test uploading a file:

```
curl --data-binary "@testfile.txt" https://admin1.hkjn.me/upload/testfile.txt
```
