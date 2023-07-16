# Mutual Trust (Knowledge Mesh)
This is a Service Composition pattern in which services aren't aware of each-other when communicating. They dispatch
tasks by broadcasting the request, other service(s) receive it and if the task is interesting they declare themselves as
being in-charge immediately. In fact task owner knows how many responses it has to collect before sending the final
response.
