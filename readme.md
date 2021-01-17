# gRPC client
This is a gRPC client for [haaukins-exercises](https://github.com/aau-network-security/haaukins-exercises) microservice which uses MongoDB to manage 
challenges on Cyber Training Platform, [Haaukins](https://github.com/aau-network-security/haaukins).
Used in our private repositories to add & update exercises which are used in the platform. 

## Requirements 

The requirement for challenge skeleton is known by team members who are contributing 
Haaukins project in terms of challenges. 
The  **__**tag**__ of the challenge skeleton** should have following format; (- all in lower case -)

__** categorytag_challengetag **__ - challenge or category tags should not include underscore sign **_**. 
It should only be exist between category and challenge tags. 

## Example Skeleton 

This is one of example skeleton which should be used in CI on Gitlab: 

```yaml

name: < Challenge Name >
tag: < categorytag_challengetag > # tag MUST BE in given format, otherwise CI will fail
instance:
  - image: < link-to-docker-image >
    dns:
      - name: < dns-record-if-any >
        type: < dns-record-type >
  - image: < link-to-docker-image >
    dns:
      - name: < dns-record-if-any >
        type: < dns-record-type >
    flags:
      - tag: < subtag-of-the-challenge >
        name: < name-of-child-challenge >
        env: < env-variable-for-flag >
        points: < challenge-points >
        category: < challenge-category-full >
        td: < challenge-description >
        od: < organization-based-desc. >
        reqs: < prerequisites-if-any >
        outcomes: <outcomes-from-the-challenge>
```
The given yaml file above is just an example of a challenge skeleton, it differs
challenge to challenge.

