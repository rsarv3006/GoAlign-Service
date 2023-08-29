#!/bin/bash

# Perform the action in a loop
for ((i=1; i<=3000; i++)); do
  # curl -X GET 'http://localhost:3000/api/task/assignedToCurrentUser' \
    # -H 'Content-Type: application/json' \
    # -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VyIjp7InVzZXJfaWQiOiI1YTQ4MzNiNS0xZTZkLTRkNTEtOGUxMS1iYzk4YWUwNTQ4NDUiLCJ1c2VyX25hbWUiOiJtZWVwIiwiZW1haWwiOiJtZWVwQHllZXQuY29tIiwiaXNfYWN0aXZlIjp0cnVlLCJpc19lbWFpbF92ZXJpZmllZCI6ZmFsc2V9LCJleHAiOjE2OTMwMDMzNTV9.CQDz5zwwyvj5FnSC4cQPqH3jHTx3HPvxu5Q-MhtbzvQ'


curl -X POST 'http://localhost:3000/api/task/' \
-H 'Content-Type: application/json' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VyIjp7InVzZXJfaWQiOiI1YTQ4MzNiNS0xZTZkLTRkNTEtOGUxMS1iYzk4YWUwNTQ4NDUiLCJ1c2VyX25hbWUiOiJtZWVwIiwiZW1haWwiOiJtZWVwQHllZXQuY29tIiwiaXNfYWN0aXZlIjp0cnVlLCJpc19lbWFpbF92ZXJpZmllZCI6ZmFsc2V9LCJleHAiOjE2OTMwMTY4MDV9.niqbA5RaA51RCVW9Yv5S5GNVJN3GzIh0ISXVhc3uTp4' \
-d '{
  "task_name": "waffles",
  "notes": "waffles are good",
  "required_completions_needed": 2,
  "team_id": "07ae863f-714d-47c8-992a-f49b95667533",
  "start_date": "2023-09-25T03:21:17Z",
  "creator_id": "5a4833b5-1e6d-4d51-8e11-bc98ae054845",
  "interval_between_windows_count": 1,
  "interval_between_windows_unit": "day(s)",
  "window_duration_count": 1, 
  "window_duration_unit": "day(s)",
  "status": "active"
}'
done
