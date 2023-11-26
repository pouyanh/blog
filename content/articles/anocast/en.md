asdas

# Anocast
This is a Communication pattern in which nodes aren't aware of each-other when communicating. They broadcast messages,
other node(s) receive it and in case of interest they declare themselves as being in-charge immediately.
In fact message publisher knows how many replies it has to collect in order to conclude.

## Realization
Suppose having a task dispatcher which broadcasts the task without having information about which unit(s)
are going to work on the task and some workers which may work on some tasks. Task dispatch and gathering
the result(s) is asynchronous. task dispatcher and workers share an array per each task accessible by
{TaskName}-{RequestID}. TaskName is equivalent to message topic in messaging systems and RequestID is a uniquely
generated id by task dispatcher.

Whenever workers receive a task, in case they realized that the task is interesting for them which means
they want to process it, they **push** to the shared array {TaskName}-{RequestID}. Value doesn't matter.

When each worker finishes working on the task it publishes its individual result so that the task dispatcher
receives it as one partial response while there may be other possible responses too. In this step the task dispatcher
*pop*s from the shared array {TaskName}-{RequestID}. If there were still remaining items on the shared array it means
there are workers which are still working on the task in fact the task dispatcher don't return the final response and
just keeps this partial response somewhere to join with other partial responses later. If there were no items left on
the shared array it means this was the partial response from last worker which means all workers have finished their
jobs, and it's time to prepare the final response by combining partial responses from all workers.
Regardless of combining logic, final response is ready to be sent.

## Boring Message
If no node is interested in the message, the message publisher waits forever to collect at least one reply which is
pointless. To handle this situation there should be an early-timeout mechanism. One solution is to check length of the
shared array {TaskName}-{RequestID} a bit (500ms~1s) after message published. If it's greater than zero it means there
are at least one node which is processing the message and message publisher has to wait for reply. If it's zero it
means no one interested in the message and the message publisher has to go on.

## Multi-staged Messaging
Some node may need to communicate with other nodes in order to reply to the received message. Don't forget that each
**node** can be a **message publisher** too. As this pattern is based on asynchronous communication nodes are free to
publish secondary messages and wait for replies confident about primary message publisher waits for it.

## Fraud Prevention
All nodes should have unique signatures and at least replies should be signed. Collector have to check if all replies to
a unique message have distinct signatures. This puts limitation on number of replies each individual node can send to a
single message. It's supposed that nodes process each message at most once.

## Compensation
Message possibly have side effects which may be desired to be withdrawn in case one node reported a problem. It can
become handled by different techniques including rollback and commit on success. One solution is to keep side effects
in an unsure state while sending replies and waiting for a supplementary message to actually apply the side effect.
Side effects not applied, expire after some time. This supplementary message is published by the same publisher of the
original message after receiving all replies and only if the message journey is successful. For instance original
message being: _FormSubmitted_ and supplementary message being: _FormSubmissionSucceed_.

## Implement using Messaging & Redis List
This pattern lets us handle synchronous communication (frontend-backend) by variable number of
asynchronous communications (service-service) which may differ due to the request payload. In the diagram below
communication between internal services are event-driven. Any async messaging platform can play role of the
_Message BUS_.

Here All 3 services are subscribed to the message being published (**Foo.Happened**) but due to the
request payload only 2 of them are interested in current request which are **Service 1** & **Service 2**.
They show their interests by pushing to the list dedicated to **ReqID=1** on Redis. **Service 2** finishes its own job
after a while and publishes the related event: **Foo.Finished**. **API Gateway** receives it and pops from the List
related to **ReqID=1**. As it still contains another element **API Gateway** realizes that there are some other
services (here **Service 1**) still processing the request. So it doesn't send the response to the **Client**.
Meanwhile, **Service 2** needs other services to assist on the procedure, so it creates a subtask by dispatching an
event (**Bar.Happened**) and waits for responses just like **API Gateway** is waiting for all services to finish their
jobs on **Foo.Happened**. **Service 3** is the only subscriber, and it shows its interest immediately after receiving
the event by pushing to the list dedicated to **ReqID=2**. After a while it finishes the process & publishes the event
which **Service 2** is waiting for: **Bar.Finished**. As it was the only service which was processing **ReqID=2** when
**Service 2** pops the only item of list related to **ReqID=2** it gets empty and so **Service 2** finds out that
all possible workers (here **Service 3**) have been done their job on the subtask. It's time to go on. When
**Service 2** finishes its own job it publishes a **Foo.Finished** event and **API Gateway** receives it. It pops from
the list related to **ReqID=2** and as it gets empty it means all workers are done, and it's time to make the final
response ready and send it back to the **Client**.

[![](https://mermaid.ink/img/pako:eNq9VU2L2zAQ_SuDDyWhScE2ezFsIB9smlJDsLs3X2R5koiNJVeSU8Ky_72jfHodh3Z76MXIevNGb-YN0qvHVYFe5Bn8WaPkOBNsrVmZSQBWWyXrMkedycM_t0rDFJiB6VagtG6zYtoKLiomLYznDhsvFzBnFn-xfTsgcXiChTBtJHZIjMawNcLkOW3jqe8CUtQ7wRH8GzhowsENHDbh8FiORm5Br_Ne8PAwgPOn76DpcDT6PJ5HMNkq_iLkmkRTe4yF3nnxpNQhdDyn2DiCZZ1vhdlcauj9UJXgjxT25SurKpRYuCSL2aPfv7RT7KhPELu_eEh5Uj-CtM4N1yKneKugyX9HSv0rK_h7VnBlhX9iFdiU6HZSn4hJBAtpUVMXiPeJKqeyiW43CN_FqUenOhtHJ4cEwQcTXEjU4ichqcNE-abyZn8v-81zG-KPRbfbPRo5gxOstnv6chQ7ytC7kyK-WE3il6oysNKqvK_YPyqeCVMxyzdogIGpc8vMy1n5hOn2ZAT9e3Nx41WT3aX1Wvxp2tPwA50PuqwL77vg1LRcCG5cCLtdcDPf7ULQ6cJ5Bu-7cMNL3rvyz3Pk_8856qrgwBxOaRhQOpmmUtKgY51W51uJYG_glahLJgq63V_dZubROSVmXkTLlXJTkHmZfKNId9Wne8m9yOoaB15dFXTs6S3wohXbGtqli5tegPj4YBzejbff50gIfA?type=png)](https://mermaid.live/edit#pako:eNq9VU2L2zAQ_SuDDyWhScE2ezFsIB9smlJDsLs3X2R5koiNJVeSU8Ky_72jfHodh3Z76MXIevNGb-YN0qvHVYFe5Bn8WaPkOBNsrVmZSQBWWyXrMkedycM_t0rDFJiB6VagtG6zYtoKLiomLYznDhsvFzBnFn-xfTsgcXiChTBtJHZIjMawNcLkOW3jqe8CUtQ7wRH8GzhowsENHDbh8FiORm5Br_Ne8PAwgPOn76DpcDT6PJ5HMNkq_iLkmkRTe4yF3nnxpNQhdDyn2DiCZZ1vhdlcauj9UJXgjxT25SurKpRYuCSL2aPfv7RT7KhPELu_eEh5Uj-CtM4N1yKneKugyX9HSv0rK_h7VnBlhX9iFdiU6HZSn4hJBAtpUVMXiPeJKqeyiW43CN_FqUenOhtHJ4cEwQcTXEjU4ichqcNE-abyZn8v-81zG-KPRbfbPRo5gxOstnv6chQ7ytC7kyK-WE3il6oysNKqvK_YPyqeCVMxyzdogIGpc8vMy1n5hOn2ZAT9e3Nx41WT3aX1Wvxp2tPwA50PuqwL77vg1LRcCG5cCLtdcDPf7ULQ6cJ5Bu-7cMNL3rvyz3Pk_8856qrgwBxOaRhQOpmmUtKgY51W51uJYG_glahLJgq63V_dZubROSVmXkTLlXJTkHmZfKNId9Wne8m9yOoaB15dFXTs6S3wohXbGtqli5tegPj4YBzejbff50gIfA)

## Join partial responses

## An example

## Alternative names
* Trust Chain/Mutual Trust: I trust you, You trust me.
* Knowledge Mesh: Workflow knowledge is decentralized. It is distributed among all workers.
