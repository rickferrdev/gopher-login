## Philosophy

This template is intentionally minimal.

It provides a small Go project structure based on the Ports and Adapters architecture, but it does not force specific infrastructure choices such as databases, queues, repositories, caches, or external clients.

The goal is to avoid repeating basic project setup while still allowing each project to evolve in its own direction.

Add ports and adapters only when your application has a real need for them.

## What this template includes

* Basic Go project layout
* Dependency injection with Uber Fx
* HTTP server setup
* Environment variable loading
* Minimal inbound and outbound module separation
* A simple place to add application services

## What this template does not include by default

* Repository interfaces
* Database adapters
* Queue adapters
* Authentication
* Observability stack
* Business-specific domain models
* Framework-specific opinions beyond the minimal HTTP setup

These pieces should be added by each project according to its own requirements.
