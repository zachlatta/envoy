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
    
4. Run the following commands to set up your environment variables -- make sure to replace the sample data with your actual values.

    ```sh
    # Twilio number you bought
    export FROM_NUMBER="+13104242222"
    # Your phone number
    export TO_NUMBER="+13102424444"
    
    # Slack API token -- get it from https://api.slack.com/docs/oauth-test-tokens
    export SLACK_TOKEN=xoxp-49598003913-50591248285-62829404320-dx13f26cc7
    # Channel to use Envoy with. Don't include the hashtag in its name
    export SLACK_CHANNEL=general

    # Twilio account SID, get it from your Twilio console
    export TWILIO_SID=AC2dbb33ca2a49fab0df2b744c236k9dc2
    # Twilio auth token, get it from your Twilio console
    export TWILIO_TOKEN=eb748ed3782949245fee2a38fe80a80b
    ```
    
5. Start the beast!

        $ go build && ./envoy

## License

See [LICENSE](LICENSE).
