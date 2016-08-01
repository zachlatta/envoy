# Envoy

![](demo.gif)

Envoy is a barebones SMS client for Slack. Enjoy Slack but hate the mobile client? This may be for you.

## How it Works

Envoy links a Slack channel with a Twilio phone number, allowing you to communicate with the given channel entirely over SMS.

Right now Envoy only allows you to communicate with a single channel, though hopefully that'll change with time.

## Setup

1. Follow the standard protocol for installing Go apps.

        $ go get -u github.com/zachlatta/envoy

2. Head over to https://twilio.com and create an account if you don't already have one.

3. Buy a phone number on Twilio and set its incoming message callback to `https://<your_envoy_url>/callback/sms` (ex. `https://envoyapp.ngrok.io/callback/sms`)

    ![](twilio_callback_setup.gif)
    
4. Run the following commands to set up your environment variables -- make sure to replace anything between `<` and `>` with proper values.

    ```sh
    export FROM_NUMBER=<Twilio number you bought (ex. +13104442222)>
    export TO_NUMBER=<your phone number (ex. +13102224444)>

    export SLACK_TOKEN=<Slack API token, get it from https://api.slack.com/docs/oauth-test-tokens>
    export SLACK_CHANNEL=<Slack channel you'd like to link with Envoy (ex. general) - no hashtag>

    export TWILIO_SID=<Twilio account SID, get on your Twilio console>
    export TWILIO_TOKEN=<Twilio auth token, get on your Twilio console>
    ```
    
5. Start the beast!

        $ go build && ./envoy

## License

See [LICENSE][LICENSE].
