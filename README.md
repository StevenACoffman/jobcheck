# jobcheck
Example kubernetes job with file health checking

Restricting what's in your runtime container to precisely what's necessary for your app is a best practice employed by Google and other tech giants that have used containers in production for many years. It improves the signal to noise of scanners (e.g. CVE) and reduces the burden of establishing provenance to just what you need.

Many applications running for long periods of time eventually transition to broken states, and cannot recover except by being restarted. Kubernetes provides liveness probes to detect and remedy such situations.

How do we monitor the health of an ephemeral job in golang, without having to embed a webserver?

A sentinel file that has it's timestamp modified on a regular basis can provide this information.

This small program is able to create and maintain such a file, and check the health of the file.