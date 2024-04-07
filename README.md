## miniblog Project

miniblog is a beginner-level Go project used to implement user registration, blog creation, and other functionalities. miniblog is both beginner-friendly and comprehensive:

- **Beginner-friendly:** It is suitable for Go developers who have just learned the basic syntax of Go and have no prior experience in project development.
- **Comprehensive:** This project originates from a large-scale online project of a first-tier enterprise. After learning miniblog, you will be equipped with the necessary knowledge to develop enterprise-level projects.

miniblog implements the following two types of functionalities:
- **User management:** It supports 7 user operations, including user registration, user login, fetching user list, fetching user details, updating user information, changing user password, and user logout.
- **Blog management:** It supports 6 blog operations, including creating a blog, fetching blog list, fetching blog details, updating blog content, deleting a blog, and batch deleting blogs.

**This project is for**

- Beginners who have just completed learning the basic syntax of Go and want to quickly learn and participate in Go language development in a company.
- Individuals who have a grasp of the basic syntax of Go but have zero experience in building Go application  and want to learn Go development in a fast and right way through practical experience.
- Those who are interested in Go application development but are still beginners or have limited experience.
- Developers who have some experience in Go application development but want to learn advanced development skills.

Courses:

- Juejin Course: [Developing Enterprise-Level Go Applications from Scratch](https://juejin.cn/book/7176608782871429175)
- Geek Time Course: [Practical Go Language Project Development](https://time.geekbang.org/column/intro/100079601?tab=intro)

## Features

- Inspired by concept of clean architecture.
- Uses most commonly used Go packages such as gorm, casbin, govalidator, jwt-go, gin, cobra, viper, pflag, zap, pprof, grpc, protobuf, and more.
- Follows a standardized directory structure based on the [project-layout](https://github.com/golang-standards/project-layout) convention.
- Includes authentication (JWT) and authorization (casbin) capabilities.
- Includes independently designed log and error packages.
- Utilizes a high-quality Makefile for project management.
- Static code analysis.
- Includes unit tests, performance tests, fuzz tests, and sample tests.
- Provides rich web functionality such as call chains, graceful shutdown, middleware, cross-origin resource sharing (CORS), and exception recovery.
  - Implements HTTP, HTTPS, and gRPC servers.
  - Supports JSON and Protobuf as data exchange formats.
- Adheres to various development standards including code style, versioning, interface design, logging, error handling, and commit conventions.
- Implements MySQL programming for data access.
- Implements business functionalities for user management and blog management.
- Follows RESTful API design guidelines.
- Provides OpenAPI 3.0/Swagger 2.0 API documentation.

> Note: The above features represent the core functionalities I selected based on senior experience as a Go developer. These functionalities also cover most of the core features in development of enterprise-level Go application . By mastering these functionalities, you will be able to develop high-quality enterprise applications after complete this project.

## Installation

```bash
$ git clone https://github.com/jxs1211/myminiblog.git
$ go work use myminiblog # If Go version > 1.18
$ cd myminiblog
$ make # Compile the source code
```

## Documentation

- [User Manual](./docs/guide/zh-CN/README.md)
- [Developer Manual](./docs/devel/zh-CN/README.md)

## Contributing

You are welcome to contribute code and star the project!

### Development Guidelines

This project follows the following development guidelines: [miniblog Project Development Guidelines](./docs/devel/zh-CN/conversions/README.md).

## License

[MIT](https://choosealicense.com/licenses/mit/)