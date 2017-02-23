# CoffeShop

This example try to demonstrates building a "real-world" web services in go.

## Description

This repo focus on building the web services in monolithic way. Although we structure the app in such way migrating from monolithic to microservices would be easy.

### Organization

This application consists of three main services. 

- User service
- Catalog Service
- Order services

To make it more "real-world", I have added a simple worker to process asynchronous tasks and simple catalog search via elastic search backend.

- __user__ - user service. Handles login, signup and reset password
- __catalog__ - catalog service. Listing and get details of an item(Book)
- __order__ - order service. Place an order.
- __worker__ - used by order to process the placed order aysnchronously
