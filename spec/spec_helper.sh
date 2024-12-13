# shellcheck shell=sh

root='/tmp/snekcheck_test'
bin='./result/bin/snekcheck'

spec_helper_precheck() {
  minimum_version '0.28.1'
}

spec_helper_configure() {
  before_each 'global_before_each_hook'
  after_each 'global_after_each_hook'
}

global_before_each_hook() {
  mkdir -p $root
}

global_after_each_hook() {
  rm -r $root
}
