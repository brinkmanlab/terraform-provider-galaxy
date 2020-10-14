resource "galaxy_user" "{{ $name }}" {
  username = "{{ $username }}"
  password = "{{ $password }}"
  email = "{{ $email }}"
}