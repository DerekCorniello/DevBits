# Title - DevBits

## Goal: Create an X and LinkedIn crossover for posting real content about your projects, semi-formally

### Outline:

Projects are the baseline of the app, you must have a project to post.

Anyone can comment, like, follow, despite their project count.

- Projects are called 'Stream's
- Posts about projects are called 'Byte's
- Comments are called 'Bit's
- (Not sure if we want this to be the names in the db or not)

### Tech Stack:

- Backend/API in Go, Elixir/Scala if need big data processing
- Frontend: ReactNative? React? Vue? VanillaJS?
- Database: PostgreSQL or MySQL
- Host: On AWS, use EC2 and API Gateway for Backend, link to RDS and Elasticache instance. Pull media from a CloudFront Instance.
