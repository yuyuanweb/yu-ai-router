/**
 * 加密工具类
 *
 * @author <a href="https://codefather.cn">编程导航学习圈</a>
 */
package com.yupi.airouter.utils;

import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import javax.crypto.Cipher;
import javax.crypto.spec.SecretKeySpec;
import java.nio.charset.StandardCharsets;
import java.util.Base64;

@Slf4j
@Component
public class EncryptionUtils {

    /**
     * 加密密钥（从配置文件读取）
     */
    @Value("${encryption.secret-key:yupi-ai-router-secret-key-256}")
    private String secretKey;

    /**
     * AES 加密
     */
    public String encrypt(String plainText) {
        try {
            // 确保密钥长度为 32 字节（256 位）
            byte[] keyBytes = getKeyBytes(secretKey);
            SecretKeySpec secretKeySpec = new SecretKeySpec(keyBytes, "AES");
            
            Cipher cipher = Cipher.getInstance("AES/ECB/PKCS5Padding");
            cipher.init(Cipher.ENCRYPT_MODE, secretKeySpec);
            
            byte[] encryptedBytes = cipher.doFinal(plainText.getBytes(StandardCharsets.UTF_8));
            return Base64.getEncoder().encodeToString(encryptedBytes);
        } catch (Exception e) {
            log.error("加密失败", e);
            throw new RuntimeException("加密失败", e);
        }
    }

    /**
     * AES 解密
     */
    public String decrypt(String encryptedText) {
        try {
            byte[] keyBytes = getKeyBytes(secretKey);
            SecretKeySpec secretKeySpec = new SecretKeySpec(keyBytes, "AES");
            
            Cipher cipher = Cipher.getInstance("AES/ECB/PKCS5Padding");
            cipher.init(Cipher.DECRYPT_MODE, secretKeySpec);
            
            byte[] decryptedBytes = cipher.doFinal(Base64.getDecoder().decode(encryptedText));
            return new String(decryptedBytes, StandardCharsets.UTF_8);
        } catch (Exception e) {
            log.error("解密失败", e);
            throw new RuntimeException("解密失败", e);
        }
    }

    /**
     * 获取固定长度的密钥字节数组（32 字节）
     */
    private byte[] getKeyBytes(String key) {
        byte[] keyBytes = new byte[32];
        byte[] sourceBytes = key.getBytes(StandardCharsets.UTF_8);
        System.arraycopy(sourceBytes, 0, keyBytes, 0, Math.min(sourceBytes.length, 32));
        return keyBytes;
    }
}
