
# Certrack

This is just a simple tls certificate tracker.

## Installation

```
git clone https://github.com/muchiri08/certrack
```

<p>Then</p>

```
cd certrack
```

<p>Execute the below command to build the project and generate the binary file</p>

```
make build
```

<p>bin/certrack will generated</p>
<p>To run you need to pass several arguments when executing the binary file: </p>

 ``--dsn``, ``--smtp-host``, ``smtp-port``, ``--smtp-user``, ``--smtp-pwd``, ``--smtp-sender`` and an optional ``--port``

 | Parameter | Description                | Example |
| :--------  | :------------------------- | :--------------------- |
| `dsn` | **Required**. Your datasource name | ``postgres://<user>:<password>@<host>/certrack``
| `smtp-host` | **Required**. Your smtp host | ``sandbox.smtp.mailtrap.io``
| `smtp-port` | **Required**. Your smtp port | ``2525``
| `smtp-user` | **Required**. Your smtp username | ``username``
| `smtp-pwd` | **Required**. Your smtp password | ``password``
| `smtp-sender` | **Required**. Your smtp sender email | ``email@example.com``
| `port` | **Optional**. Your application port | ``default is 4000``

<p>Example of the full command</p>

```
bin/certrack --dsn postgres://user:password@host/certrack --smtp-port 2525 --smtp-host sandbox.smtp.mailtrap.io --smtp-user username --smtp-pwd password --smtp-sender email@example.com
```

<p>After successfully running navigate to the host and port you have run the app on your browser eg localhost:4000</p>



