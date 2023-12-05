temp_file=$(mktemp)
tee $temp_file | ./trebuchet
awk -f trebuchet.awk < $temp_file | ./trebuchet
rm $temp_file
