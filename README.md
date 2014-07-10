# Polyglot

_Polyglot is experimental software. Please use with caution and at your own risk!_

## What is it?

Programmers use [web frameworks](http://en.wikipedia.org/wiki/Web_application_framework) to simplify the effort to develop web applications. Web frameworks reduce the overhead of work that needs to be done in most web applications. This reduces time and effort needed, improving stability and creating a consistent and maintainable system. Web application frameworks are written in a single, major programming language, like Ruby with Ruby on Rails, Python with Django, Java with Spring or Javascript with Angular.js. 

Polyglot is a web framework that *increases* the complexity in the effort to develop web applications. Unlike frameworks like Rails or Django or Express, Polyglot doesn't exist to make life easier for the programmer.

So what good is that?? 

The answer is **trade-offs**. 

As a programmer you trade-off complexity and effort for something you think is more important for the web application you're creating. In Polyglot, we are trading complexity and effort for:

1. **Performance scalability** -- Polyglot responders are distributed and independent processes that can reside anywhere on a connected network
2. **Modularity** -- Polyglot responders can be chained, each doing an individual piece of processing, encouraging reusability of code
3. **Extensibility** -- by creating an acceptor as a controller in an existing web application, you can extend the applications through Polyglot
4. **Multi-lingual development** -- Polyglot responders can be developed in multiple programming languages, **at the same time**


*Polyglot* is not for all types of web applications. You should only use Polyglot for web applications that need to be scale in a highly performant way and/or need to be incrementally developed in multiple programming languages. For example, if your web application never needs to scale beyond a single server, you're probably better off using some other single language framework. And if once you create your web application and you or anyone else never need to add new features, Polyglot is probably not for you either. 

The first three are understandable, but the fourth is quite strange, why would you want to develop a web application in multiple programming languages? There are good, practical reasons:

1. Web applications you write are systems and they change over time and can be written or maintained by different groups of people. If you're not restricted to a particular platform or language, then the chances of getting an incrementally better piece of software is higher. 
2. Also, by forcing the deliberate use of different programming languages, you are forced to separate the layers and make each component more independent and robust, being able to switch out the poor-performing responders and replacing them with higher-performing ones
3. Different responders can have different criteria for performance, ease-of-development, ease-of-maintenance or quick turnaround in development. With a single programming language you are often forced to accept a compromise. With multiple programming languages, you can choose the platform and language as what you need for that particular responder
4. Different responders can be written for specific performance gains or maintenability.


Is Polyglot for your web application? 

## Architecture

Polyglot has a very simple and basic architecture. It consists of 3 main components:

1. **Acceptor** - the acceptor is a HTTP interface that takes in a HTTP request and provides a HTTP response. The acceptor takes in a HTTP request converts it into a generic message and drops it into the message queue. Then depending on what is asked, it will return the appropriate HTTP response. The default implementation of the acceptor is in Go. To extend an existing web application, you can implement the acceptor as a controller in that web application.
2. **Messsage queue** - a queue that receives the messages that represent the HTTP request. the acceptor accepts HTTP requests and converts the requests into messages that goes into the message queue. The messages then gets picked up by the next component, the responder. The implementation of the message queue is a RabbitMQ server.
3. **Responder** - the responder is a standalone process written in any language that picks up messages from the message queue and responds to it accordingly. In most web frameworks this is usually called the controller. However unlike most web framework controllers, the responders are actual independent processes that can potentially sit in any connected remote server. Responders contain the "business logic" of your code. 

Essentially, the Polyglot architecture revolves around using a messaging queue to disassociate the processing units (responders) from the communications unit (acceptor), allowing the responders to be created in multiple languages.

###Flow

There are two basic types of flows in Polyglot.

The normal flow goes like this:

1. Client sents a HTTP request to the acceptor
2. The acceptor converts the request into JSON and adds the JSON message into the message queue, and waits for a response
3. A responder detects the message and starts processing it
4. The responder completes processing the message and adds a response message back to the message queue, with the correlation ID set to the same ID that was sent as part of the request message. The response message is an array that has 3 elements -- the HTTP status code, a map of headers, and a body.
5. The acceptor detects a message with the same correlation ID on the queue and picks it up
6. The acceptor uses the HTTP status code, response headers and the body, creates a HTTP response and sends it to the client

The chained flow goes like this:

1. Client sents a HTTP request to the acceptor
2. The acceptor converts the request into JSON and adds the JSON message into the message queue, and waits for a response
3. A responder detects the message and starts processing it
4. Once the responder completes processing, it will send a create another message on the queue for another responder to process, then waits for a response
5. This results in a chain of responders
6. Once the final responder completes processing, the results are gathered and rolled back to the first responder
7. The first responder adds a response message back to the message queue, with the correlation ID set to the same ID that was sent as part of the request message
8. The acceptor detects a message with the same correlation ID on the queue and picks it up
9. The acceptor uses the HTTP status code, response headers and the body, creates a HTTP response and sends it to the client


### Acceptor

The acceptor is a communications unit that interacts with the external world (normally a browser). The default implementation a simple web application written in Go, using the [Gin framework](http://gin-gonic.github.io/gin/) (yes, the irony of implementing a framework with another framework). The acceptor is sessionless and main task is to accept requests and push them into the message queue, then receives the response and reply to the requestor. 

You can also extend an existing application by creating a controller in that application as an acceptor. 

### Message Queue

The message queue is a generic [AMQP](http://www.amqp.org) message queue, implemented using [RabbitMQ](https://www.rabbitmq.com/). 

### Responder

Responders are processing units that can be written in any programming language that can communicate with the message queue. Responders are normally written as standalone processes that wait on the message queue. All responders essentially do the same thing, which is to process incoming requests but there are 3 types of responders that behave slightly differently:

1. First responder -- this is a responder that responds to a message from the acceptor
2. Final responder -- this is a responder that returns a message to the acceptor
3. Chained responder -- this is a responder that takes a message from a responder and creates a message to another responder

A responder can be a first and final responder if it responds to a message from the acceptor and returns a message to the same acceptor. 


## Installation and setup

My development environment in OS X Mavericks but it should work fine with most *nix-based environments. 

First, clone this repository in to a [Go workspace](http://golang.org/doc/code.html). The default acceptor is written in Go and you'll need to build it.

### Acceptor

The default polyglot acceptor is written in Go using the [Gin framework](http://gin-gonic.github.io/gin/). To install, just run:

    go build
    
This should create a program called `polyglot`. To run the default acceptor:

    ./polyglot
    
    

### Message queue

The default message queue is RabbitMQ. You can find instructions to [download and install RabbitMQ in your environment here](https://www.rabbitmq.com/download.html). A couple of notable items if you've not used RabbitMQ before:

1. The default username is 'guest' and the default password is also 'guest'
2. The web admin URL is http://localhost:15672. You can do most queue management and administration through this interface
3. The default 'guest' user can only connect via localhost, so if you're looking at deploying remote responders, please set up your [access control properly](https://www.rabbitmq.com/access-control.html)

The default implementation of Polyglot makes no additional demands on RabbitMQ but if you're looking at something more secured, I would suggest to create a dedicated exchange and implement finer grained access control on RabbitMQ for eg setting up clients to publish and get on specific queues.

Once you have installed RabbitMQ, you can start up the queue with this:

    rabbitmq-server


### Responder

The set of example responders are in the `responders` directory. To run them, you can either start them individually or you can use [Foreman](https://github.com/ddollar/foreman) like I do. 

To start the responders individually, go to the respective directories for eg the `ruby` directory and do this:

    ruby hello.rb
    
This will start up the `hello` Ruby responder. Note that you have only started 1 responder. To start up and manage a bunch of responders, open up `Procfile` in the responders directory.

    hello_ruby: ruby -C ./ruby hello.rb
    foo_ruby: ruby -C ./ruby foo.rb
    hello_php: php ./php/hello.php

You will notice that the file consists of lines of configuration that starts up the responders. Now open up the file `.foreman`.

    concurrency: 
      hello_ruby=5,
      foo_ruby=5,
      hello_php=5

As you can see, each line after `concurrency:` is a responder. The number configuration is the number of responders you want to start up. In this case, I'm starting up 5 `hello_ruby`, `foo_ruby` and `hello_php` responders each.

To start all the responders at once:

    foreman start
    
If you want to run this in production, use Foreman to export out the configuration files in Upstart or launchd etc. If you want to run this in the background without being cut off when you log out:

    nohup foreman start &


## Writing responders

### Ruby


### Python


### PHP




## Extending the default acceptor


## Extending an existing application

