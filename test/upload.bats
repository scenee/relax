#!/usr/bin/env bats

load test_helper

@test "relax upload -n crashlytics ipa" {
  export CL_TOKEN=FOO
  export CL_SECRET=BAR
  export CL_GROUP=TEST
  IPA="$(relax show development ipa)"
  run relax upload -n crashlytics "$(relax show development ipa)"
  assert_success

  [[ "${lines[@](-1)}" =~ "-ipaPath '$IPA'"  ]] \
  && [[ "${lines[@](-1)}" =~ "-groupAliases '$CL_GROUP'"  ]] \
  && [[ "${lines[@](-1)}" =~ "$CL_TOKEN"  ]] \
  && [[ "${lines[@](-1)}" =~ "$CL_SECRET"  ]]
}

@test "relax upload -n testfairy ipa" {
  export TF_API_KEY=TF_API_KEY
  export TF_METRICS=TF_METRICS
  export TF_GROUPS=TF_GROUPS
  export TF_VIDEO=TF_VIDEO
  export TF_NOTIFY=TF_NOTIFY

  IPA="$(relax show development ipa)"
  run relax upload -n testfairy  "$(relax show development ipa)"

  assert_success

  [[ "${lines[@](-1)}" =~ "-F file=@'$IPA'"  ]] \
  && [[ "${lines[@](-1)}" =~ "-F api_key='$TF_API_KEY'"  ]] \
  && [[ "${lines[@](-1)}" =~ "-F metrics='$TF_METRICS'"  ]] \
  && [[ "${lines[@](-1)}" =~ "-F notify='$TF_NOTIFY'"  ]] \
  && [[ "${lines[@](-1)}" =~ "-F testers_groups='$TF_GROUPS'"  ]]
}

