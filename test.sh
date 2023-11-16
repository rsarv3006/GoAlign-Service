#!/bin/bash

# Perform the action in a loop
for ((i=1; i<=3; i++)); do
  # curl -X GET 'http://localhost:3000/api/task/assignedToCurrentUser' \
    # -H 'Content-Type: application/json' \
    # -H 'Authorization: Bearer '


curl -X POST 'http://localhost:3000/api/v1/task/' \
-H 'Content-Type: application/json' \
-H 'Authorization: Bearer ' \
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
