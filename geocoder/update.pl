#!/usr/bin/perl

use Data::Dumper;
use HTTP::Tiny;
use JSON;
use Locale::Country;


`wget http://download.geonames.org/export/dump/allCountries.zip && unzip allCountries.zip` or die "Cannot download or unzip files";

open(my $fh, "<", "allCountries.txt") or die "Cannot open file";
my $columns = [ map { $_ =~ s/^\s+|\s+$//g; $_ =~ s/\s|\-/_/g; lc $_; } split("\t", "geonameid	name	asciiname	alternatenames	latitude	longitude	feature class	feature code	country code	cc2	admin1 code	admin2 code	admin3 code	admin4 code	population	elevation	dem	timezone	modification date") ];

my $http = HTTP::Tiny->new;
while(<$fh>) {
	my @line = map { $_ =~ s/^\s+|\s+$//g; $_; } split("\t", $_);
	my %data = map { $columns->[$_] => $line[$_] } 0..$#line;
	$data{country} = code2country(lc $data{country_code});
	my $city = { city => \%data };
	$http->post("http://localhost:9200/geo/city/", { content => to_json({ city => \%data }) });
}
close($fh);
`rm allCountries.*`; 
