A minimalist implementation of the Firefox Update Server (https://github.com/mozilla-releng/balrog).

Overview
========

There are two basic facts about Firefox releases that make serving updates for it a difficult problem:
1) Each Firefox release contains hundreds of different builds (primarily distinguished by platform and locale)
2) There are thousands of different combinations of parameters a Firefox instance may send when requesting an update

Combined, these two things make it very difficult to compactly represent the various update paths for Firefox versions.

Further adding to the complexity is the fact that humans must be able to easily and confidently read and manipulate these update paths.

High Level Overview
===================

Firefox regularly checks for updates by contacting a server with a path formatted as follows:
/update/6/{product}/{version}/{buildID}/{buildTarget}/{locale}/{channel}/{osVersion}/{systemCapabilities}/{distribution}/{distVersion}/update.xml

Upon receiving a request Gothmog does the following:
* Splits and parses the fields above (some may contain multiple distinct values)
* Uses its rules engine to determine which one single release the request should receive
* Generates an appropriate response based on the metadata of the determined release

Notably, an actual update package is not returned - merely a pointer to one. For example:
```
<updates>
    <update type="minor" displayVersion="56.0 Beta 3" appVersion="56.0" platformVersion="56.0" buildID="20170815141045" detailsURL="https://www.mozilla.org/en-GB/firefox/56.0/releasenotes/">
        <patch type="complete" URL="http://download.mozilla.org/?product=firefox-56.0b3-complete&amp;os=linux64&amp;lang=en-GB" hashFunction="sha512" hashValue="d0df415e4de8830f086f82bcfd78159906167e6c01797b083879d456e4f0c702b1bd8876283473672df1c38abf953a98336eeadc41b39d645fbc020a3a22bd84" size="53408965"/>
    </update>
</updates>
```

Components
==========

Rule Definitions
----------------
Rules are how we define how incoming requests are mapped to particular releases. Whenever possible, an update request should simply receive the latest version of Firefox, but there are a number of reasons why this cannot always be the case ("watershed" updates, where the current version cannot apply the latest update directly, but must update to an inbetween version first, platform deprecation, and other reasons).

Rules contain a number of properties. The majority of these simply correspond to fields from the request path above (these are generally a 1-1 mapping, except for `systemCapabilities`, which is split into multiple properties). The other properties are:
* release_mapping - The name of the release that update requests should receive if they resolve to this rule
* priority - The priority of this rule, relative to the other defined rules. This is used to ensure rules are parsed deterministically.

Release Models
--------------
Each release is represented as a JSON document containing the necessary metadata to construct responses for all platforms and locales built for a particular version. Its contents are too lengthy to show here, but you can see an example at TODO.

Rules Engine
------------
The rules engine is responsible for taking an incoming update request, evaluating it against the rules in the system, and determining which single rule is the best match for the request. At a very high level, the algorithm is simple:
1) Throw away all of the rules who have properties defined that don't match the equivalent field in the incoming request. For example, if the rule defines a channel of "release", but the incoming request has a channel of "beta" - the rule would be thrown away.
2) Choose the rule with the highest priority.

Unsurprisingly, step 1 is where the real complexity lies. In practice, whether or not a proprety matches an incoming field is often more than a simple string match, and special rules exist for most of them. This special matching is the key to how gothmog is able to achieve its goal of compactly representing updating paths while keeping them understandable.

First off, any property of a rule that is null will always match any incoming value in an update request. This feature alone reduces the number of necessary rules in the system by at least an order of magnitude.

Now, let's look at a few properties and how they are matched.

The `version` property supports an optional </>/<=/>= prefix to allow matching more than just a single, specific version. Eg: "<72.0". This is regularly used in so-called "watershed" updates -- where an update to an intermediate version of Firefox is required before updating to the latest one.

The `channel` property supports an optional `*` glob at the end of a specified channel. This is commonly used to ensure production channels are set-up equivalently to internal testing channels, reducing the risk that testing misses a bug that is hit in production.

Response Generator
------------------
