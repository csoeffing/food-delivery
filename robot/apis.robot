*** Settings ***

Documentation   Functional tests for Agents Webservice
Library         OperatingSystem
Library         RequestsLibrary
Library         Collections


*** Variables ***
${BASE_URL}     http://localhost:8134

*** Keywords ***
Get With API Token
    [Arguments]    ${apiToken}    ${endpoint}
    ${headers}=    Create Dictionary    Authorization=Bearer ${apiToken}
    ${response}=   GET    ${BASE_URL}${endpoint}    headers=${headers}    expected_status=anything
    RETURN    ${response}

*** Test Cases ***

User Login
    ${loginData}=       Create Dictionary   user_name=test1     password=test
    ${response}=        POST                ${BASE_URL}/api/user/login     json=${loginData}
	${payload}=         Set Variable        ${response.json()}
    ${apiToken}=        Set Variable        ${payload["access_token"]}

    #Log To Console      ${apiToken}
	
    ${userResp}=        Get With API Token      ${apiToken}            /api/user/1
	${userData}=        Set Variable            ${userResp.json()}

    #Log To Console      ${userData}

	Should Be Equal As Integers     ${userData["user_id"]}        ${1}
	Should Be Equal As Strings      ${userData["user_name"]}        test1
	