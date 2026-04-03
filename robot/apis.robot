*** Settings ***

Documentation   Functional tests for Agents Webservice
Library         OperatingSystem
Library         RequestsLibrary
Library         Collections
Suite Setup     Wait Until Keyword Succeeds    10 seconds    1 second      Connect To Server

*** Variables ***
${BASE_URL}     http://localhost:8134

*** Keywords ***

Connect To Server
    ${loginData}=       Create Dictionary   userName=gwash     password=test
    ${response}=        POST                ${BASE_URL}/api/user/login     json=${loginData}

Get With API Token
    [Arguments]    ${apiToken}    ${endpoint}
    ${headers}=    Create Dictionary    Authorization=Bearer ${apiToken}
    ${response}=   GET    ${BASE_URL}${endpoint}    headers=${headers}    expected_status=anything
    RETURN    ${response}

*** Test Cases ***

User Login
    ${loginData}=       Create Dictionary   userName=gwash     password=test
    ${response}=        POST                ${BASE_URL}/api/user/login     json=${loginData}
	${payload}=         Set Variable        ${response.json()}
    ${apiToken}=        Set Variable        ${payload["access_token"]}

    #Log To Console      ${apiToken}
	
    ${userResp}=        Get With API Token      ${apiToken}            /api/user/1
	${userData}=        Set Variable            ${userResp.json()}

    #Log To Console      ${userData}

	Should Be Equal As Integers     ${userData["id"]}        ${1}
	Should Be Equal As Strings      ${userData["userName"]}        gwash

Duplicate Users
    ${user1}=     Create Dictionary
    ...           firstName=George
    ...           lastName=Washington
    ...           userName=gwash
    ...           password=test
    ...           email=gwash@gmail.com
    ...           phone=800-555-1235

	${user2}=     Create Dictionary
    ...           firstName=George
    ...           lastName=Washington
    ...           userName=gwash
    ...           password=test
    ...           email=gwash2@gmail.com
    ...           phone=800-555-1235

	${user3}=     Create Dictionary
    ...           firstName=George
    ...           lastName=Washington
    ...           userName=gwash2
    ...           password=test
    ...           email=gwash2@gmail.com
    ...           phone=800-555-1235
    
    #Run Keyword And Expect Error  *400 Bad*     POST    ${BASE_URL}/api/user/signup     json=${user1}
	${response1}=        POST            ${BASE_URL}/api/user/signup     json=${user1}    expected_status=400
    ${response1Code}=    Set Variable    ${response1.status_code}
    ${response1Json}=    Set Variable    ${response1.json()}

    Should Be Equal As Integers          ${response1.status_code}            400
	Should Contain                       ${response1.json()["message"]}      Email address already exists

	${response2}=        POST            ${BASE_URL}/api/user/signup     json=${user2}    expected_status=400

    Should Be Equal As Integers          ${response2.status_code}            400
	Should Contain                       ${response2.json()["message"]}      Username already exists

	${response3}=        POST            ${BASE_URL}/api/user/signup     json=${user3}    expected_status=400

    Should Be Equal As Integers          ${response3.status_code}            400
	Should Contain                       ${response3.json()["message"]}      Phone number already exists