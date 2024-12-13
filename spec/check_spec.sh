# shellcheck shell=sh

Describe "check"
  Context "with an empty, valid directory"
    It "succeeds"
      When call "$bin" "$root"
      The status should be success
    End
  End

  Context "with a single invalid file"
    create_invalid_file() { touch "$root"/InVaLiD; }
    BeforeEach "create_invalid_file"

    It "fails"
      When call "$bin" "$root"
      The status should be failure
    End

    It "does not modify the file"
      When call "$bin" "$root"
      The file "$root"/InVaLiD should be exist
    End
  End
End
