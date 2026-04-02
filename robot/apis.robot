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
    ${loginData}=       Create Dictionary   userName=test1     password=test
    ${response}=        POST                ${BASE_URL}/api/user/login     json=${loginData}

Get With API Token
    [Arguments]    ${apiToken}    ${endpoint}
    ${headers}=    Create Dictionary    Authorization=Bearer ${apiToken}
    ${response}=   GET    ${BASE_URL}${endpoint}    headers=${headers}    expected_status=anything
    RETURN    ${response}

*** Test Cases ***

User Login
    ${loginData}=       Create Dictionary   userName=test1     password=test
    ${response}=        POST                ${BASE_URL}/api/user/login     json=${loginData}
	${payload}=         Set Variable        ${response.json()}
    ${apiToken}=        Set Variable        ${payload["access_token"]}

    #Log To Console      ${apiToken}
	
    ${userResp}=        Get With API Token      ${apiToken}            /api/user/1
	${userData}=        Set Variable            ${userResp.json()}

    #Log To Console      ${userData}

	Should Be Equal As Integers     ${userData["id"]}        ${1}
	Should Be Equal As Strings      ${userData["userName"]}        test1

Duplicate Users
    ${user1}=     Create Dictionary
    ...           firstName=George
    ...           lastName=Washington
    ...           userName=test1
    ...           password=test
    ...           email=user2@gmail.com
    ...           phone=800-555-1235
    
    Run Keyword And Expect Error  *400*     POST    ${BASE_URL}/api/user/signup     json=${user1}