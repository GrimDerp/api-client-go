#/usr/bin/env/sh

set -o nounset
set -o errexit

api-client-go help 2>&1 | grep 'Reads functions'
api-client-go help 2>&1 | grep readsets
api-client-go help 2>&1 | grep 'Datasets'

openssl aes-256-ctr -d -a -in oauth_for_travis.pem.enc -out oauth_for_travis.pem -pass "pass:$KEY_PASSWORD"

api-client-go readsets search --use-oauth=oauth_for_travis.json --dataset-ids=376902546192 | grep ChhDSkRta1luOENoQ2otTmJNb2UzaWxiWUI=

api-client-go readsets search --use-oauth=oauth_for_travis.json --dataset-ids=376902546192 --page-token=ChhDSkRta1luOENoQ2otTmJNb2UzaWxiWUI= | grep "Page Token"

api-client-go readsets search --use-oauth=oauth_for_travis.json --dataset-ids=376902546192 HG00310 | grep CJDmkYn8ChCj-NbMoe3ilbYB

api-client-go readsets get --use-oauth=oauth_for_travis.json CJDmkYn8ChCj-NbMoe3ilbYB | grep HG00310

api-client-go reads search --use-oauth=oauth_for_travis.json --readset-ids=CJDmkYn8ChCj-NbMoe3ilbYB | grep "Page Token"

api-client-go datasets list --use-oauth=oauth_for_travis.json | grep " 376902546192,"

./api-client-go datasets get --use-oauth=oauth_for_travis.json | grep "ID: 376902546192,"
./api-client-go datasets get --use-oauth=oauth_for_travis.json --dataset-id=376902546192 | grep "Is Public: true"
