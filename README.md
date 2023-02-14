# Workouts-GO

This application is for managing a workouts platform. It allows trainers set availability
for when they can accommodate a training. Users can be the everyday people (attendee) that can see the available
times and schedule or cancel a training or the user can be the trainer. The application
itself is structured to have HTTP APIs for the users, trainers and trainings.
Each service manages its respective domain.

## Thought Process for Application Design

For each trainer we want to be able to mark the hours a trainer is available or unavailable and to
view all the available hours. An hour is available only if a trainer lists it as a possible hour and
there is no scheduled training for that hour.

The trainings API allows us to get all the trainings, create a new training with notes and the time,
cancel a training using the UUID of the existing training, reschedule a training with the UUID,
approve a rescheduled training or reject it. Each training data contains information about the UUID,
user ID and user's name, the time, notes, if it can be canceled

For every user we can get the current user, role (trainer or user), name and balance (credits used for scheduling training)

We define GRPC for RPC to let services talk to each other rather than over HTTP. GRPC provides many advantages such as
low payloads, streaming, HTTP/2 > HTTP/1, low latency high throughput, RTC.

For the users service, other services need to get balance and update balance
