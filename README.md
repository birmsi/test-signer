# test-signer
toggl - Unattended Programming Test


You can find the endpoints at - internal\signatures\api

Had some doubts about:
- what to do with the jwt if i should use it all or just some of it. Ended using just an "userid" :)
- The questions that were sent on the sign. It seemed that they werent relevan to the generation of the sign hash nor to be returned so i accepted an array of questions and "discarded" it.
- If they were to be used i probably would have created an array of QuestionAnswer to map easily the answer of each question.


To persist the data i've used a postgres DB hosted on neon.tech - the create script is the following:
```
CREATE TABLE IF NOT EXISTS user_signatures (
    user_id VARCHAR(255) NOT NULL,
    signature BYTEA NOT NULL,
    answers TEXT[] NOT NULL,
    hash_timestamp TIMESTAMP NOT NULL,
    PRIMARY KEY (user_id, signature)
); 
```
