General settings are specified at the root level of cfg.yaml. They include general configuration that applies to the Hooker application.

![settings](img/hooker-settings.png)

Key | Description                                                                                                                                                                          | Possible Values | Example Value
--- |--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------| --- | ---
*khulnasoft-server*| Khulnasoft Platform URL. This is used for some of the integrations to will include a link to the Khulnasoft UI                                                                                   | Khulnasoft Platform valid URL | https://server.my.khulnasoft
*db-verify-interval*| Specify time interval (in hours) for Hooker to perform database cleanup jobs. Default: 1 hour                                                                                        | any integer value  | 1
*max-db-size*| The maximum size of Hooker database (in B, KB, MB or GB). Once reached to size limit, Hooker will delete old cached messages. If empty then Hooker database will have unlimited size | any integer value with a unit siffux | 200kb, 1000 MB, 1Gb

