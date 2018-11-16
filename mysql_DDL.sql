DROP FUNCTION IF EXISTS getParList;
delimiter //
CREATE FUNCTION `getParList`(rootId INT)
RETURNS varchar(1000)
BEGIN
    DECLARE sTemp VARCHAR(1000);
    DECLARE sTempPar VARCHAR(1000);
    SET sTemp = '';
    SET sTempPar = rootId;

    #循环递归
    WHILE sTempPar is not null DO
        #判断是否是第一个，不加的话第一个会为空
        IF sTemp != '' THEN
            SET sTemp = concat(sTemp,',',sTempPar);
        ELSE
            SET sTemp = sTempPar;
        END IF;
        SELECT group_concat(uid) INTO sTempPar FROM kc_competition where uid<>cid and FIND_IN_SET(cid,sTempPar)>0;
    END WHILE;

RETURN sTemp;
END
//
