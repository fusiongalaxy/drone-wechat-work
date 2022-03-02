# drone-wechat-work

Wechat Work robot plugin for Drone

### Usage

```yaml
  - name: notify
    image: fusiongalaxy/drone-wechat-work
    settings:
      url: https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=XXX-XXXX-XXX-XXXXX
      msgtype: markdown
      content: |
        {{if eq .Status "success" }}
        #### 🎉 ${DRONE_REPO} 构建成功
        > Commit: [${DRONE_COMMIT_MESSAGE}](${DRONE_COMMIT_LINK})
        > Author: ${DRONE_COMMIT_AUTHOR}
        > [点击查看](${DRONE_BUILD_LINK})
        {{else}}
        #### ❌ ${DRONE_REPO} 构建失败
        > Commit: [${DRONE_COMMIT_MESSAGE}](${DRONE_COMMIT_LINK})
        > Author: ${DRONE_COMMIT_AUTHOR}
        > 请立即修复!!!
        > [点击查看](${DRONE_BUILD_LINK})
        {{end}}
    when:
      status:
        - failure
        - success
```



## Options

| option | type | required | default | description |
| --- | --- | --- | --- | --- |
| url | string | Yes | none | The full address of webhook: https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxxxxx |
| type | string | No | text | message type，support (text,markdown) |
| content | string | Yes | none |  Message content, text or markdown or json string |
| mentioned_list | string | No | none | Comma separated `userid`s to be mentioned, for example: "ponyma,jackma" or "@all" for all users |
| mentioned_mobile_list | string | No | none | Comma separated mobile phone numbers to be mentioned, for example: "13812345678,13898754321" or "@all" for all users |
