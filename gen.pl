#!/usr/bin/perl
use strict;
use warnings;

use Getopt::Long;

GetOptions("o|output=s" => \(my $output),
           "ip=s"       => \(my $ip),
           "domain=s"   => \(my $domain),
           "h|help"     => \(my $help),
           "v|verbose"  => \(my $verbose));

sub usage($) {
    my $code = shift;
    my $msg = <<EOF;
[options]  generate tsl cert.
    -h
    --help        Print help message

    --domain      你的域名
    --ip          你的ip地址

    -v
    --verbose     Print debug log
EOF
    print($msg);
    exit($code);
}

if($help) {
    usage(0);
}

if( (not $ip) || (not $domain)){
    usage(1);
}

END {
    system("rm -f server.csr");
    system("rm -f ca.srl");
    system("rm -rf _openssl.cnf");
};
sub mk_root_crt(){
    system("openssl req -new -newkey rsa:2048 -days 365 -nodes -x509 -subj \"/C=CN/ST=HuBei/L=WuHan/O=BBFE\" -keyout ca.key  -out ca.crt")
}

sub mk_openssh_cnf(@) {
    my $ip = shift;
    my $domain = shift;
    $ENV{"IP"} = $ip;
    $ENV{"DOMAIN"} = $domain;
    system("envsubst < openssl.cnf > _openssl.cnf");
}

sub mk_cert_request() {
    system("openssl genrsa -out server.key 2048");
    system("openssl req -new -key server.key -out server.csr -config _openssl.cnf");
}
sub mk_sign(){
    system("openssl x509 -req -days 120 -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -extensions v3_req -extfile _openssl.cnf -out server.crt");
    system("openssl x509 -in server.crt -noout -text");
}
mk_root_crt();
mk_openssh_cnf($ip,$domain);
mk_cert_request();
mk_sign();
