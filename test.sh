#!/bin/bash

# Perform the action in a loop
for ((i=1; i<=3; i++)); do
  # curl -X GET 'http://localhost:3000/api/task/assignedToCurrentUser' \
    # -H 'Content-Type: application/json' \
    # -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VyIjp7InVzZXJfaWQiOiI1YTQ4MzNiNS0xZTZkLTRkNTEtOGUxMS1iYzk4YWUwNTQ4NDUiLCJ1c2VyX25hbWUiOiJtZWVwIiwiZW1haWwiOiJtZWVwQHllZXQuY29tIiwiaXNfYWN0aXZlIjp0cnVlLCJpc19lbWFpbF92ZXJpZmllZCI6ZmFsc2V9LCJleHAiOjE2OTMwMDMzNTV9.CQDz5zwwyvj5FnSC4cQPqH3jHTx3HPvxu5Q-MhtbzvQ'


curl -X POST 'http://localhost:3000/api/v1/task/' \
-H 'Content-Type: application/json' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VyIjp7InVzZXJfaWQiOiI5NjU3YTU2ZC05MDRhLTRiNDYtYjdiNy1lMjY5YWZmMDAzZDEiLCJ1c2VybmFtZSI6Im1lZXAiLCJlbWFpbCI6Im1lZXBAeWVldC5jb20iLCJpc19hY3RpdmUiOnRydWUsImlzX2VtYWlsX3ZlcmlmaWVkIjpmYWxzZSwiY3JlYXRlZF9hdCI6IjIwMjMtMDktMjVUMjI6MDg6NDguNDgzMjkxWiJ9LCJleHAiOjE2OTYzMDI2MDF9.VeKIGa0cEyiwa3IIEOnbsRFEiin-FiB26_3FYR8aeP0' \
-d '{
  "task_name": "waffles",
  "notes": "waffles are good",
  "required_completions_needed": 2,
  "team_id": "ed60d73d-3b82-4929-b290-eda258b4cb00",
  "start_date": "2023-11-25T03:21:17Z",
  "interval_between_windows":{"interval_count":1,"interval_unit":"day(s)"},
  "window_duration":{"interval_count":1,"interval_unit":"day(s)"},
  "status": "active",
  "assigned_user_id": "c0356d25-e386-484e-b825-a5270808ecb4"

}'
done
