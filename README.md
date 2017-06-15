# maild
mail http service

## Usage

```
Usage of ./maild:
  -from string
        from (default=user)
  -pass string
        smtp pass
  -port int
        smtp port (default 25)
  -smtpaddr string
        smtp address
  -user string
        smtp user

```

## Demo

```
./maild -smtpaddr your_smtpaddress -pass your_password -user example@example.com

curl -s -F receiver="example@example.com" -F subject="hello" -F body="body text" localhost:3001
```