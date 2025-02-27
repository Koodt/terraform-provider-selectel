---
layout: "selectel"
page_title: "Selectel: selectel_domains_record_v1"
sidebar_current: "docs-selectel-resource-domains-record-v1"
description: |-
  Manages a V1 record resource within Selectel Domains API Service.
---

# selectel\_domains\_record\_v1

Manages a V1 record resource within Selectel Domains API Service.

## Example usage

```hcl
resource "selectel_domains_domain_v1" "domain_1" {
  name = "testdomain.xyz"
}

resource "selectel_domains_record_v1" "cname_record_1" {
  domain_id = selectel_domains_domain_v1.domain_1.id
  name = "cname.testdomain.xyz"
  type = "CNAME"
  content = "origin.com"
  ttl = 60
}

resource "selectel_domains_record_v1" "ns_record_1" {
  domain_id = selectel_domains_domain_v1.domain_1.id
  name = "testdomain.xyz"
  type = "NS"
  content = "ns5.selectel.org"
  ttl = 86400
}


resource "selectel_domains_record_v1" "a_record_1" {
  domain_id = selectel_domains_domain_v1.domain_1.id
  name = "a.testdomain.xyz"
  type = "A"
  content = "127.0.0.1"
  ttl = 60
}

resource "selectel_domains_record_v1" "aaaa_record_1" {
  domain_id = selectel_domains_domain_v1.domain_1.id
  name = "aaaa.testdomain.xyz"
  type = "AAAA"
  content = "2400:cb00:2049:1::a29f:1804"
  ttl = 60
}

resource "selectel_domains_record_v1" "txt_record_1" {
  domain_id = selectel_domains_domain_v1.domain_1.id
  name = "txt.testdomain.xyz"
  type = "TXT"
  content = "hello, world!"
  ttl = 60
}

resource "selectel_domains_record_v1" "mx_record_1" {
  domain_id = selectel_domains_domain_v1.domain_1.id
  name = "mx.testdomain.xyz"
  type = "MX"
  content = "mail.example.org"
  ttl = 60
  priority = 10
}

resource "selectel_domains_record_v1" "srv_record_1" {
  domain_id = selectel_domains_domain_v1.domain_1.id
  name = "srv.testdomain.xyz"
  type = "SRV"
  target = "backupbox.example.com"
  ttl = 120
  priority = 10
  weight = 20
  port = 100
}

resource "selectel_domains_record_v1" "caa_record_1" {
  domain_id = selectel_domains_domain_v1.main_domain.id
  name = format("caa.%s", selectel_domains_domain_v1.main_domain.name)
  type = "CAA"
  tag = "issue"
  flag = 128
  value = "letsencrypt.com"
  ttl = 60
}

resource "selectel_domains_record_v1" "sshfp_record_1" {
  domain_id = selectel_domains_domain_v1.main_domain.id
  name = format("%s", selectel_domains_domain_v1.main_domain.name)
  type = "SSHFP"
  algorithm = 1
  fingerprint_type = 1
  fingerprint = "01AA"
  ttl = 60
}
resource "selectel_domains_record_v1" "alias_record_1" {
  domain_id = selectel_domains_domain_v1.main_domain.id
  name = format("subc.%s", selectel_domains_domain_v1.main_domain.name)
  type = "ALIAS"
  content = format("%s", selectel_domains_domain_v1.main_domain.name)
  ttl = 60
}
```

## Argument Reference

The following arguments are supported:

* `domain_id` - (Required) Represents an identifier of the associated domain.
 Changing this creates a new domain record.

* `name` - (Required) Represents a name of the domain record.

* `type` - (Required) Represents a type of the record.
 Possible values: A, AAAA, TXT, CNAME, NS, SOA, MX, SRV.

* `ttl` - (Required) Represents a time-to-live for the record.
 Must be the value between 60 and 604800.

* `content` - (Optional) Represents a content of the record.
 Absent for SRV records.

* `email` - (Optional) Represents an email of the domain's admin.
 For SOA records only.

* `priority` - (Optional) Represents priority of the records preferences.
 Lower value means more preferred. MX/SRV records only.

* `weight` - (Optional) Represents a relative weight for records with the same priority,
 higher value means higher chance of getting picked.
 For SRV records only.

* `port` - (Optional) Represents TCP or UDP port on which the service is to be found.
 For SRV records only.

* `target` - (Optional) Represents a canonical hostname of the machine providing the service.
 For SRV records only.

* `tag` - (Optional) Represents the identifier of the property represented by the record.
 For CAA records only.

* `flag` - (Optional) Represents the critical flag, that has a specific meaning per RFC.
 For CAA records only.

* `value` - (Optional) Represents a value associated with the tag.
 For CAA records only.

* `algorithm` - (Optional) Represents the algorithm of the public key.
 For SSHFP records only.

* `fingerprint_type` - (Optional) Represents an algorithm used to hash the public key.
 For SSHFP records only.

* `fingerprint` - (Optional) Represents a hexadecimal hash result, as text.
 For SSHFP records only.

## Attributes Reference

The following attributes are exported:

* `content` - Represents a content of the record.

* `email` - Represents an email of the domain's admin.

* `priority` - Represents priority of the records preferences.

* `weight` - Represents a relative weight for records with the same priority,
 higher value means higher chance of getting picked.

* `port` - Represents TCP or UDP port on which the service is to be found.

* `target` - Represents a canonical hostname of the machine providing the service.

* `tag` - Represents the identifier of the property represented by the record.

* `flag` - Represents the critical flag, that has a specific meaning per RFC.

* `value` - Represents a value associated with the tag.

* `algorithm` - Represents th algorithm of the public key.

* `fingerprint_type` - Represents an algorithm used to hash the public key.

* `fingerprint` - Represents a hexadecimal hash result, as text.

## Import

Domain records can be imported using a combined ID in the following format: ``<domain_id>/<record_id>``

```shell
$ env SEL_TOKEN=SELECTEL_API_TOKEN terraform import selectel_domains_record_v1.record_1 45623/123
```
